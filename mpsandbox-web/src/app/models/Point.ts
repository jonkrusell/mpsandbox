export class Point {
    x: number
    y: number
    GUID: string
    Name: string

    constructor(x, y: number) {
        this.x = x;
        this.y = y;
        this.GUID = Guid.newGuid();
        this.Name = "";
    }
}

export class Guid {
    static newGuid() {
        return 'xxxxxxxx-xxxx-4xxx-yxxx-xxxxxxxxxxxx'.replace(/[xy]/g, function(c) {
            var r = Math.random()*16|0, v = c == 'x' ? r : (r&0x3|0x8);
            return v.toString(16);
        });
    }
}