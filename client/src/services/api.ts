import ky, {type Options} from 'ky';

const defaultInstance = ky.create({
  prefixUrl: import.meta.env.VITE_BACKEND_URL,
});

interface HttpRequestOptions extends Options {
  noAuthHeader?: boolean;
}

const authorizationMiddleware = async (request: Request): Promise<void> => {
  const options = request as Request & HttpRequestOptions;
  if (options.noAuthHeader) {
    return;
  }

  const token = localStorage.getItem("token");
  if (token) {
    request.headers.set('Authorization', `Bearer ${token}`);
  }
};

const authInstance = defaultInstance.extend({
  hooks: {
    beforeRequest: [authorizationMiddleware],
  },
})

export {defaultInstance, authInstance};
