import { fileLoader } from 'nest-typed-config';
import { RootConfig } from 'src/config/root';

let config: RootConfig;

export const loadConfig = () => {
  if (config) {
    return config;
  }

  config = fileLoader({
    absolutePath: process.env.CONFIG_FILE ?? 'config.yml',
  })() as RootConfig;
  return config;
};
