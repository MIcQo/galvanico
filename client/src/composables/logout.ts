import {useRouter} from "vue-router";
import {removeToken} from "@/services/auth.ts";
import {useAlert} from "@/stores/alert.ts";
import {useI18n} from "vue-i18n";

export function useLogout() {
  const router = useRouter();
  const alert = useAlert();
  const {t} = useI18n();

  const logout = () => {
    removeToken()
    alert.open(t("auth.alert.logout"));
    router.push({name: "auth.login"})
  }

  return {logout};
}
