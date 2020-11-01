import User from "./User";
import Location from "./Location";
import Space from "./Space";
import Ajax from "../util/Ajax";

export default class Search {
    users: User[]
    locations: Location[]
    spaces: Space[]

    constructor() {
        this.users = [];
        this.locations = [];
        this.spaces = [];
    }

    deserialize(input: any): void {
        if (input.users) {
            this.users = input.users.map(user => {
                let e = new User();
                e.deserialize(user);
                return e;
            });
        }
        if (input.locations) {
            this.locations = input.locations.map(location => {
                let e = new Location();
                e.deserialize(location);
                return e;
            });
        }
        if (input.spaces) {
            this.spaces = input.spaces.map(space => {
                let e = new Space();
                e.deserialize(space);
                return e;
            });
        }
    }

    static async search(keyword: string): Promise<Search> {
        return Ajax.get("/search/" + keyword).then(result => {
            let e: Search = new Search();
            e.deserialize(result.json);
            return e;
        });
    }
}