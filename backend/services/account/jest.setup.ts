import * as dotenv from 'dotenv';
import * as fs from 'fs';
import * as path from 'path';

const envTestPath = path.resolve(process.cwd(), '.env.test');
const envPath = fs.existsSync(envTestPath)
  ? envTestPath
  : path.resolve(process.cwd(), '.env');

dotenv.config({ path: envPath, quiet: true });
