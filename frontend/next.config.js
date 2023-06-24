// next.config.js
const path = require('path');
const dotenv = require('dotenv');

// Resolve the parent directory path
const parentDir = path.join(__dirname, '..');

// Load the .env file from the parent directory
const envConfig = dotenv.config({ path: path.join(parentDir, '.env') }).parsed;

module.exports = {
  env: {
    BACKEND_BASE_URL: envConfig.BACKEND_HOST + ":" + envConfig.BACKEND_PORT, 
  },
};
