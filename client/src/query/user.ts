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
  return useQuery({queryKey: ['currentUser'], queryFn: fetchCurrentUser})
}

async function fetchCurrentUser(): Promise<UserResponse> {
  return defaultInstance.get('api/user').json();
}
