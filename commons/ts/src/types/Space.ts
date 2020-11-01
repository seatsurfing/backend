import { Entity } from "./Entity";
import Ajax from "../util/Ajax";
import Location from "./Location";

export default class Space extends Entity {
    name: string;
    x: number;
    y: number;
    width: number;
    height: number;
    rotation: number;
    available: boolean;
    locationId: string;
    location: Location;

    constructor() {
        super();
        this.name = "";
        this.x = 0;
        this.y = 0;
        this.width = 0;
        this.height = 0;
        this.rotation = 0;
        this.available = false;
        this.locationId = "";
        this.location = new Location();
    }

    serialize(): Object {
        return Object.assign(super.serialize(), {
            "name": this.name,
            "x": this.x,
            "y": this.y,
            "width": this.width,
            "height": this.height,
            "rotation": this.rotation,
        });
    }

    deserialize(input: any): void {
        super.deserialize(input);
        this.name = input.name;
        this.locationId = input.locationId;
        this.x = input.x;
        this.y = input.y;
        this.width = input.width;
        this.height = input.height;
        this.rotation = input.rotation;
        if (input.available) {
            this.available = input.available;
        }
        if (input.location) {
            this.location.deserialize(input.location);
        }
    }

    getBackendUrl(): string {
        return "/location/"+this.locationId+"/space/";
    }

    async save(): Promise<Space> {
        return Ajax.saveEntity(this, this.getBackendUrl()).then(() => this);
    }

    async delete(): Promise<void> {
        return Ajax.delete(this.getBackendUrl() + this.id).then(() => undefined);
    }

    static async get(locationId: string, id: string): Promise<Space> {
        return Ajax.get("/location/"+locationId+"/space/" + id).then(result => {
            let e: Space = new Space();
            e.deserialize(result.json);
            return e;
        });
    }

    static async list(locationId: string): Promise<Space[]> {
        return Ajax.get("/location/"+locationId+"/space/").then(result => {
            let list: Space[] = [];
            (result.json as []).forEach(item => {
                let e: Space = new Space();
                e.deserialize(item);
                list.push(e);
            });
            return list;
        });
    }

    static async listAvailability(locationId: string, enter: Date, leave: Date): Promise<Space[]> {
        let payload = {
            enter: enter.toISOString(),
            leave: leave.toISOString()
        };
        return Ajax.postData("/location/"+locationId+"/space/availability", payload).then(result => {
            let list: Space[] = [];
            (result.json as []).forEach(item => {
                let e: Space = new Space();
                e.deserialize(item);
                list.push(e);
            });
            return list;
        });
    }
}