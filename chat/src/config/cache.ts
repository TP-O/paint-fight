import { SecondsTime } from '@enum/time';
import { CacheModuleOptions } from '@nestjs/cache-manager';
import fsStore from 'cache-manager-fs';

export const CacheConfig = Object.freeze<CacheModuleOptions>({
  store: fsStore, // Support writing data into disk
  options: {
    path: 'diskcache',
    ttl: 10 * SecondsTime.Miniute,
    maxsize: 1000 * 1000 * 1000, // bytes
    zip: true,
    preventfill: true,
  },
});
