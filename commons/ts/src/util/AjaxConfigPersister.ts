import AjaxCredentials from "./AjaxCredentials";

export default interface AjaxConfigPersister {
    persistRefreshTokenInLocalStorage(c: AjaxCredentials): Promise<void>
    readRefreshTokenFromLocalStorage(): Promise<AjaxCredentials>
    updateCredentialsSessionStorage(c: AjaxCredentials): Promise<void>
    readCredentialsFromSessionStorage(): Promise<AjaxCredentials>
    deleteCredentialsFromSessionStorage(): Promise<void>
}
