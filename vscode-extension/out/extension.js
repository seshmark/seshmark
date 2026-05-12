"use strict";
var __createBinding = (this && this.__createBinding) || (Object.create ? (function(o, m, k, k2) {
    if (k2 === undefined) k2 = k;
    var desc = Object.getOwnPropertyDescriptor(m, k);
    if (!desc || ("get" in desc ? !m.__esModule : desc.writable || desc.configurable)) {
      desc = { enumerable: true, get: function() { return m[k]; } };
    }
    Object.defineProperty(o, k2, desc);
}) : (function(o, m, k, k2) {
    if (k2 === undefined) k2 = k;
    o[k2] = m[k];
}));
var __setModuleDefault = (this && this.__setModuleDefault) || (Object.create ? (function(o, v) {
    Object.defineProperty(o, "default", { enumerable: true, value: v });
}) : function(o, v) {
    o["default"] = v;
});
var __importStar = (this && this.__importStar) || (function () {
    var ownKeys = function(o) {
        ownKeys = Object.getOwnPropertyNames || function (o) {
            var ar = [];
            for (var k in o) if (Object.prototype.hasOwnProperty.call(o, k)) ar[ar.length] = k;
            return ar;
        };
        return ownKeys(o);
    };
    return function (mod) {
        if (mod && mod.__esModule) return mod;
        var result = {};
        if (mod != null) for (var k = ownKeys(mod), i = 0; i < k.length; i++) if (k[i] !== "default") __createBinding(result, mod, k[i]);
        __setModuleDefault(result, mod);
        return result;
    };
})();
Object.defineProperty(exports, "__esModule", { value: true });
exports.activate = activate;
exports.deactivate = deactivate;
const vscode = __importStar(require("vscode"));
const child_process_1 = require("child_process");
const util_1 = require("util");
const execAsync = (0, util_1.promisify)(child_process_1.exec);
let enabled = true;
let disposables = [];
let statusBarItem;
function activate(context) {
    const config = vscode.workspace.getConfiguration('seshmark');
    enabled = config.get('enabled', true);
    // Status bar
    statusBarItem = vscode.window.createStatusBarItem(vscode.StatusBarAlignment.Right, 100);
    statusBarItem.command = 'seshmark.toggle';
    context.subscriptions.push(statusBarItem);
    updateStatusBar();
    // Commands
    context.subscriptions.push(vscode.commands.registerCommand('seshmark.blame', runBlameCommand), vscode.commands.registerCommand('seshmark.query', runQueryCommand), vscode.commands.registerCommand('seshmark.who', runWhoCommand), vscode.commands.registerCommand('seshmark.toggle', toggleDecorations));
    // Update decorations on file change / save
    context.subscriptions.push(vscode.window.onDidChangeActiveTextEditor(() => updateDecorations()), vscode.workspace.onDidSaveTextDocument(() => updateDecorations()), vscode.workspace.onDidOpenTextDocument(() => updateDecorations()), vscode.workspace.onDidChangeConfiguration(e => {
        if (e.affectsConfiguration('seshmark')) {
            const cfg = vscode.workspace.getConfiguration('seshmark');
            enabled = cfg.get('enabled', true);
            updateDecorations();
            updateStatusBar();
        }
    }));
    // Initial update
    updateDecorations();
}
function deactivate() {
    disposables.forEach(d => d.dispose());
}
function updateStatusBar() {
    if (enabled) {
        statusBarItem.text = '$(record-keys) Seshmark';
        statusBarItem.tooltip = 'Seshmark: AI attribution enabled. Click to toggle.';
        statusBarItem.backgroundColor = undefined;
    }
    else {
        statusBarItem.text = '$(circle-slash) Seshmark';
        statusBarItem.tooltip = 'Seshmark: AI attribution disabled. Click to toggle.';
        statusBarItem.backgroundColor = new vscode.ThemeColor('statusBarItem.warningBackground');
    }
    statusBarItem.show();
}
function toggleDecorations() {
    enabled = !enabled;
    const config = vscode.workspace.getConfiguration('seshmark');
    config.update('enabled', enabled, true);
    updateDecorations();
    updateStatusBar();
}
async function updateDecorations() {
    // Clear previous decorations
    disposables.forEach(d => d.dispose());
    disposables = [];
    const editor = vscode.window.activeTextEditor;
    if (!editor || !enabled)
        return;
    const filePath = editor.document.uri.fsPath;
    const workspaceFolder = vscode.workspace.getWorkspaceFolder(editor.document.uri);
    if (!workspaceFolder)
        return;
    try {
        const attributions = await getAttributions(filePath, workspaceFolder.uri.fsPath);
        applyDecorations(editor, attributions);
    }
    catch (err) {
        console.error('Seshmark error:', err);
    }
}
async function getAttributions(filePath, cwd) {
    // Run git blame --porcelain
    const { stdout } = await execAsync(`git blame --porcelain "${filePath}"`, { cwd });
    const lines = stdout.split('\n');
    const attributions = [];
    let current = {};
    for (const line of lines) {
        if (line.match(/^[0-9a-f]{40}/)) {
            // New commit block
            if (current.line !== undefined) {
                attributions.push(current);
            }
            const parts = line.split(' ');
            current = { commit: parts[0] };
        }
        else if (line.startsWith('author ')) {
            current.author = line.substring(7);
        }
        else if (line.startsWith('author-time ')) {
            const ts = parseInt(line.substring(12));
            current.date = new Date(ts * 1000).toISOString();
        }
        else if (line.startsWith('\t')) {
            // Tab-indicated line content
            current.line = attributions.length + 1;
        }
    }
    if (current.line !== undefined) {
        attributions.push(current);
    }
    // Fetch AI metadata for unique commits
    const uniqueCommits = [...new Set(attributions.map(a => a.commit))];
    const commitMeta = new Map();
    for (const commit of uniqueCommits) {
        try {
            const { stdout: body } = await execAsync(`git log -1 --format=%B ${commit}`, { cwd });
            const meta = parseTrailers(body);
            commitMeta.set(commit, meta);
        }
        catch {
            commitMeta.set(commit, { session: null, agent: null, model: null });
        }
    }
    // Merge metadata
    for (const attr of attributions) {
        const meta = commitMeta.get(attr.commit);
        if (meta) {
            attr.aiSession = meta.session;
            attr.aiAgent = meta.agent;
            attr.aiModel = meta.model;
            attr.isHuman = !meta.session && !meta.agent;
        }
    }
    return attributions;
}
function parseTrailers(body) {
    let session = null;
    let agent = null;
    let model = null;
    for (const line of body.split('\n')) {
        if (line.startsWith('AI-Session: ')) {
            session = line.substring(12);
        }
        else if (line.startsWith('AI-Agent: ')) {
            agent = line.substring(10);
        }
        else if (line.startsWith('AI-Model: ')) {
            model = line.substring(10);
        }
    }
    return { session, agent, model };
}
function applyDecorations(editor, attributions) {
    const config = vscode.workspace.getConfiguration('seshmark');
    const showGutter = config.get('showAgentInGutter', true);
    const colors = config.get('colors', {});
    // Group by agent
    const byAgent = new Map();
    for (const attr of attributions) {
        const key = attr.isHuman ? 'human' : (attr.aiAgent || 'unknown');
        if (!byAgent.has(key))
            byAgent.set(key, []);
        byAgent.get(key).push(attr);
    }
    for (const [agent, lines] of byAgent) {
        const color = colors[agent || 'unknown'] || '#6E7681';
        if (agent && agent !== 'human') {
            // Create decoration type
            const decorationType = vscode.window.createTextEditorDecorationType({
                gutterIconSize: 'contain',
                after: showGutter ? {
                    contentText: `[${agent.substring(0, 8)}]`,
                    color: color,
                    fontStyle: 'italic',
                    margin: '0 0 0 1em'
                } : undefined,
                overviewRulerColor: color,
                overviewRulerLane: vscode.OverviewRulerLane.Right,
            });
            // Create decoration options with hover messages
            const options = lines.map(attr => {
                const hoverMsg = new vscode.MarkdownString();
                hoverMsg.appendMarkdown(`**AI-Session:** ${attr.aiSession || 'N/A'}\n\n`);
                hoverMsg.appendMarkdown(`**Agent:** ${attr.aiAgent || 'N/A'}\n\n`);
                if (attr.aiModel) {
                    hoverMsg.appendMarkdown(`**Model:** ${attr.aiModel}\n\n`);
                }
                hoverMsg.appendMarkdown(`**Commit:** \`${attr.commit.substring(0, 7)}\`\n\n`);
                hoverMsg.appendMarkdown(`**Author:** ${attr.author || 'N/A'}`);
                hoverMsg.isTrusted = true;
                return {
                    range: new vscode.Range(attr.line - 1, 0, attr.line - 1, 0),
                    hoverMessage: hoverMsg
                };
            });
            editor.setDecorations(decorationType, options);
            disposables.push(decorationType);
        }
    }
}
// Command implementations
async function runBlameCommand() {
    const editor = vscode.window.activeTextEditor;
    if (!editor) {
        vscode.window.showInformationMessage('Open a file first.');
        return;
    }
    const terminal = vscode.window.createTerminal('Seshmark');
    terminal.sendText(`git agentblame "${editor.document.uri.fsPath}"`);
    terminal.show();
}
async function runQueryCommand() {
    const agent = await vscode.window.showInputBox({ prompt: 'Agent name (e.g., cursor, claude)' });
    if (!agent)
        return;
    const terminal = vscode.window.createTerminal('Seshmark');
    terminal.sendText(`seshmark query --agent ${agent}`);
    terminal.show();
}
async function runWhoCommand() {
    const editor = vscode.window.activeTextEditor;
    if (!editor) {
        vscode.window.showInformationMessage('Open a file first.');
        return;
    }
    const position = editor.selection.active;
    const line = position.line + 1;
    const workspaceFolder = vscode.workspace.getWorkspaceFolder(editor.document.uri);
    if (!workspaceFolder)
        return;
    try {
        const attributions = await getAttributions(editor.document.uri.fsPath, workspaceFolder.uri.fsPath);
        const attr = attributions.find(a => a.line === line);
        if (!attr || attr.isHuman) {
            vscode.window.showInformationMessage(`Line ${line}: Human-written.`);
            return;
        }
        const msg = `Line ${line}: ${attr.aiAgent || 'AI'} | ${attr.aiSession || 'no session'} | ${attr.commit.substring(0, 7)}`;
        vscode.window.showInformationMessage(msg);
    }
    catch (err) {
        vscode.window.showErrorMessage(`Seshmark error: ${err}`);
    }
}
//# sourceMappingURL=extension.js.map