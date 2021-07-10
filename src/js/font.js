
export function getPath(font, str, size) {
    const path = font.getPath(str, 0, 0, 72);
    const builder = new PathBuilder(path);
    //builder.draw( document.querySelector('#fish'));
    return builder.polygons().paths;
}

class PathBuilder {

    fontTrace;

    /** @type {number} number of points to put into corners */
    resolution = 6;

    paths = false;
    holes = false;

    path = null;
    hole = null;

    constructor(fontTrace) {
        this.fontTrace = fontTrace;
    }
    
    /** Draws the font on the given canvas */
    draw(canvas) {
        if (canvas.getContext) {
            const polygons = this.polygons();
            const ctx = canvas.getContext('2d', { });
            const offsetX = 72;
            const offsetY = 72;
            ctx.strokeStyle = 'black';
	        ctx.lineWidth = 1;
            for (let poly of polygons.paths) {
                //const poly = polygons.paths[0];
                ctx.moveTo(poly[0][0] + offsetX, poly[0][1] + offsetY);
                for(let i = 1; i < poly.length; i++) {
                    const coord = poly[i];
                    const x = coord[0] + offsetX;
                    const y = coord[1] + offsetY;
                    ctx.lineTo(x,y);
                }
                ctx.stroke();
            }
        }
    }

    build() {
        this.paths = [];
        this.holes = [];

        for (let i = 0; i < this.fontTrace.commands.length; i++) {
            let cmd = this.fontTrace.commands[i];
    
            //console.log(cmd);
            switch(cmd.type) {
                case 'M':
                    if (this.path == null) {
                        this.beginPath();
                        this.path.push([ cmd.x, cmd.y ]);
                    } else {
                        this.beginHole();
                        this.hole.push([ cmd.x, cmd.y ]);
                    }
                    break;
                case 'L':
                    this.push([ cmd.x, cmd.y ]);
                    break;
                case 'C':
                    const bs = this.coord();
                    for (let p = 0; p < this.resolution; p++) {
                        const t = p / this.resolution;
                        const bezier = getBezierXY(t, bs[0], bs[1], cmd.x1, cmd.y1, cmd.x2, cmd.y2, cmd.x, cmd.y);
                        this.push(bezier);
                    }
                    break;
                case 'Q':
                    const qs = this.coord();
                    for (let p = 0; p < this.resolution; p++) {
                        const t = p / this.resolution;
                        const quadratic = getQuadraticXY(t, qs[0], qs[1], cmd.x1, cmd.y1, cmd.x, cmd.y);
                        this.push(quadratic);
                    }
                    break;
                case 'Z':
                    this.closePath();
                    this.closeHole();
                    break;
    
            }
        }

    }

    /** Returns the paths and holes */
    polygons() {
        if (this.paths == false || this.holes == false) this.build();
        return {
            points: this.paths.reduce((tally, path) => tally + path.length, 0) + this.holes.reduce((tally, path) => tally + path.length, 0),
            paths: this.paths,
            holes: this.holes
        }
    }

    // Gets the last coord
    coord() {
        if (this.hole != null) {
            return this.hole[this.hole.length - 1];
        } else {
            return this.path[this.path.length - 1];
        }
    }

    // Pushes the coord to the stack
    push(coord) {
        if (this.hole != null) {
            this.hole.push(coord);
            return this.hole;
        } else {
            this.path.push(coord);
            return this.path;
        }
    }

    beginPath() {
        if (this.path != null) 
            this.closePath();
        this.path = [];
    }
    closePath() {
        if (this.path == null) 
            return false;

        this.path.push(this.path[0]);
        this.paths.push(this.path);
        this.path = null;
    }


    beginHole() {
        if (this.hole != null) 
            this.closeHole();
        this.hole = [];
    }
    closeHole() {
        if (this.hole == null) 
            return false;
        this.hole.push(this.hole[0]);
        this.holes.push(this.hole);
        this.hole = null;
    }
}

// gets a point of the bezier for the given coord http://www.independent-software.com/determining-coordinates-on-a-html-canvas-bezier-curve.html
function getBezierXY(t, sx, sy, cp1x, cp1y, cp2x, cp2y, ex, ey) {
    return [
        Math.pow(1-t,3) * sx + 3 * t * Math.pow(1 - t, 2) * cp1x + 3 * t * t * (1 - t) * cp2x + t * t * t * ex,
        Math.pow(1-t,3) * sy + 3 * t * Math.pow(1 - t, 2) * cp1y + 3 * t * t * (1 - t) * cp2y + t * t * t * ey
    ];
}
function getQuadraticXY(t, sx, sy, cp1x, cp1y, ex, ey) {
    return [
        (1-t) * (1-t) * sx + 2 * (1-t) * t * cp1x + t * t * ex,
        (1-t) * (1-t) * sy + 2 * (1-t) * t * cp1y + t * t * ey
    ];
}