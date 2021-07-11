import opentype from 'opentype.js';
import earcut from 'earcut';

export class Font {

    /** @type {opentype.Font} open type font */
    font;

    /**
     * Creates a new font object
     * @param {opentype.Font} font 
     */
    constructor(font) {
        this.font = font;
    }

    /** Builds a polygon for the given string */
    polygon(str, size) {
        console.groupCollapsed("Font Polygon");
        const { root, polygons } = this.builder.polygon(str, size);
        console.groupEnd("Font Polygon");
        return polygons;
    }

    /** Builds a mesh for the givne string */
    mesh(str, size) {
        console.groupCollapsed("Font Triangulate");
        const result = this.builder.triangulate(str, size);
        console.groupEnd("Font Triangulate");
        return result;
    }

    /** The builder */
    get builder() {
        return new FontBuilder(this.font);
    }
}


const MAX_BEZIER_STEPS = 10;
const BEZIER_STEP_SIZE = 3.0;
const CURVE_RESOLUTION = 8;
const EPSILON = 1e-6;
class FontBuilder {
    constructor(font) {
        this.font = font;
    }

    /** Gereates a polygon */
    polygon(str, fontSize) {
        console.log('[font]', '[polygons]', 'call', str, fontSize);

        const path      = this.font.getPath(str, 0, 0, fontSize);
        const polygons  = [];       // List of all polygons
        let polygon     = null;     // Current polygon

        console.log('[font]', '[polygons]', 'path:', path);

        /** Closes teh current polygon */
        function closePolygon() {
            if (polygon == null) return false;
            polygon.close();
            polygons.push(polygon);
            polygon = null;
            return true;
        }

        // Iterate over all the commands
        for (let { x, y, x1, y1, x2, y2, type } of path.commands) {
            switch(type) {
                case 'M':
                    if (polygon != null)
                        closePolygon();

                    polygon = new Polygon();
                    polygon.push({x, y});
                    break;
                case 'L':
                    polygon.push({x, y});
                    break;
                case 'C':
                    polygon.pushBezier({x: x1, y: y1}, {x: x2, y: y2}, {x, y});
                    break;
                case 'Q':
                    polygon.pushQuadratic({x: x1, y: y1}, {x, y});
                    break;
                case 'Z':
                    closePolygon();
                    break;

            }
        }

        // Push the final polygon
        if (polygon != null) 
            closePolygon();
        
        // Sort by area descending
        polygons.sort((a, b) => Math.abs(b.area) - Math.abs(a.area));
            
        // classify contours to find holes and their 'parents'
        const root = [];
        for (let i = 0; i < polygons.length; ++i) {
            let parent = null;
            for (let j = i - 1; j >= 0; --j) {
                // a contour is a hole if it is inside its parent and has different winding
                if (polygons[j].inside(polygons[i].points[0]) && polygons[i].area * polygons[j].area < 0) {
                    parent = polygons[j];
                    break;
                }
            }
            if (parent) {
                parent.children.push(polygons[i]);
            } else {
                root.push(polygons[i]);
            }
        }
       
        console.log('[font]', '[polygons]', root, polygons);
        return { root, polygons };
    }

    /** Triangulates */
    triangulate(str, fontSize) {
        console.log('[font]', '[triangulate]', str, fontSize);
        const { root, polygons } = this.polygon(str, fontSize);
        console.log('[font]', '[triangulate]', root, polygons);

        const pointCount = polygons.reduce((sum, p) => sum + p.points.length, 0);
        let vertexCount = 0;

        const verticies = new Float32Array(pointCount * 2);
        const indices = [];

        function buildVerts(poly) {
            console.log('[font]', '[triangulate]', 'buildVerts', poly);
            
            // construct input for earcut
            const coords = [];
            const holes = [];

            // Add our initial points
            for (let point of poly.points) 
                coords.push(point.x, point.y);

            for (let child of poly.children) {
                // Calculate the chidlren of the child
                for (var c2 of child.children)
                buildVerts(c2);
                
                // Add the child points
                holes.push(coords.length / 2);
                for (let point of child.points) 
                    coords.push(point.x, point.y);
            }

            // Add the data
            verticies.set(coords, vertexCount * 2);
            const tris = earcut(coords, holes);
            for (let i of tris) indices.push(i + vertexCount);
            vertexCount += coords.length / 2;
        }

        // Create the mesh
        for(const poly of root)
            buildVerts(poly);

        console.log('[font]', '[triangulate]', verticies.length, indices.length);
        return {
            verticies:  verticies,
            indices:    new Uint16Array(indices)
        }
    }
}

class Polygon {

    points = [];
    children = [];
    _area = 0.0;

    constructor() {}

    /** Pushes a coordinate */
    push(coord) {
        this.points.push(coord);
        this._area = false;
    }

    /** Pushes a bezier curve */
    pushBezier(cp1, cp2, endPoint) {

        /*
        const p0 = this.lastPoint;
        const dist = distance(p0, cp1) + distance(cp1, cp2) + distance(cp2, e);
        const steps = Math.max(2, Math.min(MAX_BEZIER_STEPS, dist / BEZIER_STEP_SIZE));
        for (let i = 1; i <= steps; ++i) {
          const t = i / steps;
          const a = lerp(lerp(p0, cp1, t), lerp(cp1, cp2, t), t);
          const b = lerp(lerp(cp1, cp2, t), lerp(cp2, e, t), t);
          this.push(lerp(a, b, t));
        }
        */

        
        const s = this.lastPoint;
        const sx = s.x, sy = s.y;
        const ex = endPoint.x, ey = endPoint.y;
        const cp1x = cp1.x, cp1y = cp1.y;
        const cp2x = cp2.x, cp2y = cp2.y;

        for (let p = 0; p < CURVE_RESOLUTION; p++) {
            const t = p / CURVE_RESOLUTION;
            this.push({
                x: Math.pow(1-t,3) * sx + 3 * t * Math.pow(1 - t, 2) * cp1x + 3 * t * t * (1 - t) * cp2x + t * t * t * ex,
                y: Math.pow(1-t,3) * sy + 3 * t * Math.pow(1 - t, 2) * cp1y + 3 * t * t * (1 - t) * cp2y + t * t * t * ey
            });
        }
        
    }

    /** Pushes a quadratic coordinate */
    pushQuadratic(cp1, endPoint) {
        /*
        const p0 = this.lastPoint;
        const dist = distance(p0, cp1) + distance(cp1, e);
        const steps = Math.max(2, Math.min(MAX_BEZIER_STEPS, dist / BEZIER_STEP_SIZE));
        for (let i = 1; i <= steps; ++i) {
          const t = i / steps;
          this.push(lerp(lerp(p0, cp1, t), lerp(cp1, e, t), t));
        }
        */
        
        const s = this.lastPoint;
        const sx = s.x, sy = s.y;
        const ex = endPoint.x, ey = endPoint.y;
        const cp1x = cp1.x, cp1y = cp1.y;

        for (let p = 0; p < CURVE_RESOLUTION; p++) {
            const t = p / CURVE_RESOLUTION;
            this.push({
                x: (1-t) * (1-t) * sx + 2 * (1-t) * t * cp1x + t * t * ex,
                y: (1-t) * (1-t) * sy + 2 * (1-t) * t * cp1y + t * t * ey
            });
        }
    }

    /** Closes the polygon */
    close() {
        //original implementation:
        //let cur = this.points[this.points.length - 1];
        //this.points.forEach(next => {
        //this.area += 0.5 * cross(cur, next);
        //  cur = next;
        //});
    }

    /** Checks if the polygon is within another polygon
     * @param {Polygon} p
     */
    inside(p) {
        let count = 0, cur = this.lastPoint;
        this.points.forEach(next => {
            const p0 = (cur.y < next.y ? cur : next);
            const p1 = (cur.y < next.y ? next : cur);
            if (p0.y < p.y + EPSILON && p1.y > p.y + EPSILON) {
                if ((p1.x - p0.x) * (p.y - p0.y) > (p.x - p0.x) * (p1.y - p0.y)) {
                    count += 1;
                }
            }
            cur = next;
        });
        return (count % 2) !== 0;
    }

    /** Checks the last coordinate */
    get lastPoint() {
        return this.points[this.points.length - 1];
    }

    /** Gets the area of the shape */
    get area() {
        if (this._area === false) {
            this._area = this.points.reduce((a, next) => {
                const [ tally, cur ] = a;
                return [ tally + (0.5 * (cur.x * next.y - cur.y * next.x)), next ];
            }, [ 0, this.lastPoint ])[0];
        }
        return this._area;
    }
}

/** Distance between two points */
function distance(p1, p2) {
    const dx = p1.x - p2.x, dy = p1.y - p2.y;
    return Math.sqrt(dx * dx + dy * dy);
}

/** Linearly interperlates between two points over t */
function lerp(p1, p2, t) {
    return {x: (1 - t) * p1.x + t * p2.x, y: (1 - t) * p1.y + t * p2.y};
}