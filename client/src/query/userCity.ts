import {useQuery} from '@tanstack/vue-query'
import {defaultInstance} from "@/services/api.ts";

enum Building {
  city_hall
}

export interface CityBuilding {
  Building: Building,
  Level: number,
  Position: number
}

export interface UserCitiesResponse {
  ID: string,
  Name: string,
  PositionX: number,
  PositionY: number,
  Buildings: CityBuilding[],
}

export function getUserCities() {
  return useQuery({
    queryKey: ['userCities'],
    queryFn: fetchUserCities
  })
}

async function fetchUserCities(): Promise<UserCitiesResponse[]> {
  return defaultInstance.get('api/user/city').json();
}
