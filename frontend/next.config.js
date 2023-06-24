// // next.config.js
// const path = require('path');
// const dotenv = require('dotenv');

// // Resolve the parent directory path
// const parentDir = path.join(__dirname, '..');

// // Load the .env file from the parent directory
// const envConfig = dotenv.config({ path: path.join(parentDir, '.env') }).parsed;

// module.exports = {
//   env: {
//     BACKEND_BASE_URL: envConfig.BACKEND_HOST + ":" + envConfig.BACKEND_PORT, 
//   },
// };

// next.config.js
const path = require('path');
const dotenv = require('dotenv');

// Resolve the parent directory path
const parentDir = path.join(__dirname, '..');

// Try to load the .env file from the parent directory
const dotenvConfig = dotenv.config({ path: path.join(parentDir, '.env') });

// Use the parsed dotenv config if available, otherwise use process.env
const envConfig = dotenvConfig.error ? process.env : dotenvConfig.parsed;

module.exports = {
  env: {
    BACKEND_BASE_URL: envConfig.BACKEND_HOST + ":" + envConfig.BACKEND_PORT,
  },
};

