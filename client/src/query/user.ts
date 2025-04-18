import {useQuery} from '@tanstack/vue-query'
import {defaultInstance} from "@/services/api.ts";

export interface UserResponse {
  user: User
}

export interface User {
  Username: string;
  Email: string;
}

export function getCurrentUser() {
  const {data: currentUser} = useQuery({queryKey: ['currentUser'], queryFn: fetchCurrentUser})

  return {currentUser}
}

async function fetchCurrentUser(): Promise<UserResponse> {
  return defaultInstance.get('api/user').json();
}
