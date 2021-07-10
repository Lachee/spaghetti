
const loaderUtils   = require('loader-utils');
const mime          = require('mime-types');
const path          = require('path');

const MIME_TYPES = {
    '.glsl': 'text/plain',
    '.frag': 'text/plain',
    '.vert': 'text/plain',
};

function getMimeType(resourcePath) {
    const ext = path.extname(resourcePath);
    if (MIME_TYPES[ext] != undefined) 
        return MIME_TYPES[ext];
    const mimeType = mime.contentType(ext);
    if (mimeType !== false) return mimeType;
    return 'application/octet-stream';
}

//Encodes the content
function encodeContent(content, mimeType) {
    return `data:${mimeType};base64,` + content.toString('base64');
}
//Encodes the url
function encodeUri(url, mimeType) {
    return `data:${mimeType};uri,${url}`;
}

module.exports = function(content) {
    const options = loaderUtils.getOptions(this);
    
    const { resourcePath } = this;
    const context = options.context || this.rootContext;
    
    // Get the mime type
    const mimeType = getMimeType(resourcePath);

    // Generate the URL and return it. Prefix the URL with // so we know   
    const name = options.name || `[path][name].[ext]`; 
    const url = loaderUtils.interpolateName(this, name, {
        context,
        content,
    });
    
    // Determine if it should be embedded or not
    let embed = false;
    if (options.embed) {
        if (typeof options.embed === 'function') {
            embed = options.embed(url, mimeType, context);
        } else {
            embed = options.embed;
        }
    }

    // If we should embed, do so
    if (embed) {
        let buffer = content;
        if (typeof buffer === 'string')
            buffer = Buffer.from(content);

        const encoded = encodeContent(buffer, mimeType);
        if (encoded !== undefined) {
            return 'module.exports='+JSON.stringify(encoded);
        }
    }

    // Just return the uri
    return 'module.exports='+JSON.stringify(encodeUri(url, mimeType));
}

module.exports.raw = true;