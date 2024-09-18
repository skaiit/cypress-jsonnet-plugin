const os = require('os');
const fs = require('fs');
const { execSync } = require('child_process');
const path = require('path');
const pino = require('pino');

const jsonnetModulePath = path.join('./node_modules/cy-jsonnet/jsonnet');

let commandToRun = "";

// Configure pino logger
const logger = pino({
  level: 'info',
  prettyPrint: { colorize: true }
});

/**
 * Build code to generate the jsonnet binary with custom extensions
 * Set the commandToRun so other functions can run it for jsonnet interpretation
 */
function buildPlugin() {
  try {
    logger.info('Generating cy-jsonnet binary');
    execSync(`go build -C ${jsonnetModulePath}`);
    let runner = path.join(jsonnetModulePath, 'cy-jsonnet');
    commandToRun = `${runner}`;
    // If Windows OS, setup accordingly 
    if (os.platform().includes('win32')) {
      runner = path.join(jsonnetModulePath, 'cy-jsonnet.exe');
      commandToRun = `cmd /c ${runner}`;
    } else {
      logger.info('Setting the shell script permission to be executable');
      execSync(`chmod +x ${runner}`);
    }
    logger.info(`Command to run = ${commandToRun}`);
  } catch (err) {
    logger.error('An error occurred during jsonnet extension build');
    logger.error(err);
    process.exit(1);
  }
}