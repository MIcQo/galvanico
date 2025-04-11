import ky, {type Options} from 'ky';
import {getToken} from "@/services/auth.ts";

const authorizationMiddleware = async (request: Request): Promise<void> => {
  const options = request as Request & HttpRequestOptions;
  if (options.noAuthHeader) {
    return;
  }

  const token = getToken();
  if (token) {
    request.headers.set('Authorization', `Bearer ${token}`);
  }
};

const defaultInstance = ky.create({
  prefixUrl: import.meta.env.VITE_BACKEND_URL,
  hooks: {
    beforeRequest: [authorizationMiddleware],
  },
});

interface HttpRequestOptions extends Options {
  noAuthHeader?: boolean;
}

export {defaultInstance, type HttpRequestOptions};
