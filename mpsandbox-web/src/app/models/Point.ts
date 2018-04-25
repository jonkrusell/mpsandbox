export class Point {
    x: number
    y: number

    constructor(x, y: number) {
        this.x = x;
        this.y = y;
    }
}

export class NPC {
    Point: Point
    GUID: string
    Name: string

    constructor(x, y: number) {
        this.Point = new Point(x, y);
        this.GUID = Guid.newGuid();
        this.Name = "";
    }
}

export class Player {
    Point: Point
    GUID: string
    Name: string
    Health: number

    constructor(x, y: number) {
        this.Point = new Point(x, y);
        this.GUID = Guid.newGuid();
        this.Name = "";
        this.Health = 100;
    }
}

export class Projectile {
    Point: Point
    XSpeed: number
    YSpeed: number
}

export class Shield {
    Point: Point
    Health: number
}

export class Guid {
    static newGuid() {
        return 'xxxxxxxx-xxxx-4xxx-yxxx-xxxxxxxxxxxx'.replace(/[xy]/g, function(c) {
            var r = Math.random()*16|0, v = c == 'x' ? r : (r&0x3|0x8);
            return v.toString(16);
        });
    }
}