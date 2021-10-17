import { User } from "flexspace-commons";

export default class RuntimeConfig {
    static EMBEDDED: boolean = false;

    static async setLoginDetails(context: any): Promise<void> {
        return User.getSelf().then(user => {
            context.setDetails(user.email);
        });
    }
}