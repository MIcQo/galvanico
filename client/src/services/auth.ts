const isAuthenticated = (): boolean => {
  return !!getToken();
}

const getToken = (): string | null => {
  return sessionStorage.getItem("token");
}

const setToken = (token: string): void => {
  sessionStorage.setItem("token", token);
}

const removeToken = (): void => {
  sessionStorage.removeItem("token");
}

export {isAuthenticated, getToken, setToken, removeToken};
