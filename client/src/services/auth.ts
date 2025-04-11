const isAuthenticated = (): boolean => {
  return !!getToken();
}

const getToken = (): ?string => {
  return sessionStorage.getItem("token");
}

const setToken = (token: string): void => {
  sessionStorage.setItem("token", token);
}

export {isAuthenticated, getToken, setToken};
