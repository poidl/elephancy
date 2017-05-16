import http = require('http')

export function make_request(options: any): Promise<http.IncomingMessage> {
    let req = http.request(options)
    req.end()
    return new Promise<http.IncomingMessage>((resolve, reject) => {
        req.on('response', function (message: http.IncomingMessage) {
            resolve(message)
        })
        // thrown an error if e.g. connection is refused (e.g. wrong port number)
        req.on('error', function (err: Error) {
            console.log(err)
            throw err
        })
    })
}

export function get_code_body(message: http.IncomingMessage): Promise<{ code: number, body: string }> {
    let body: string = ''
    let err: Error
    message.on('error', function (err: Error) {
        err = err
    });
    message.on('data', function (chunk) {
        body += chunk;
    });
    let code = message.statusCode
    return new Promise<{ code: number, body: string }>((resolve, reject) => {
        if (err) {
            console.log(err)
            throw (err)
        }
        message.on('end', function () {
            if (code >= 200 && code <= 299) {
                resolve({ code: code, body: body })
            } else {
                reject({ code: code, body: body });
            }
        });
    })
}




