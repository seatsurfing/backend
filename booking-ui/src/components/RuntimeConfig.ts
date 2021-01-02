import { Ajax, User } from "flexspace-commons";

export default class RuntimeConfig {
    static EMBEDDED: boolean = false;

    static async setLoginDetails(token: string, context: any): Promise<void> {
        window.sessionStorage.setItem("jwt", token);
        Ajax.JWT = token;
        return User.getSelf().then(user => {
            context.setDetails(token != null ? token : "", user.email);
        });
    }
}