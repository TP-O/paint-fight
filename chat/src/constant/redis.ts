import { CACHE } from './cache';

export const REDIS = Object.freeze({
  NAMESAPCE: {
    PLAYER_ID_TO_SOCKET_ID: `chat:${CACHE.KEY.PLAYER_ID_TO_SOCKET_ID}`,
  },
});
