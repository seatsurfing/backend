import { Entity } from "../types/Entity";
interface AjaxResult {
  json: any
  status: number
  objectId: string
}

export default class Ajax {
  static URL: string = "";
  static JWT: string = "";

  static getBackendUrl(): string {
    let url = Ajax.URL.trim();
    if (url.endsWith("/")) {
      url = url.substring(0, url.length-1);
    }
    return url;
  }

  static async query(method: string, url: string, data?: any): Promise<AjaxResult> {
    url = Ajax.getBackendUrl() + url;
    let headers = new Headers();
    if (Ajax.JWT) {
      headers.append("Authorization", "Bearer " + Ajax.JWT)
    }
    if (data && !(data instanceof File)) {
      headers.append("Content-Type", "application/json");
    }
    let options: RequestInit = {
      method: method,
      mode: "cors",
      cache: "no-cache",
      credentials: "same-origin",
      headers: headers
    };
    if (data) {
      if (data instanceof File) {
        options.body = data;
      } else {
        options.body = JSON.stringify(data);
      }
    }
    return new Promise<AjaxResult>(function (resolve, reject) {
      fetch(url, options).then((response) => {
        if (response.status >= 200 && response.status <= 299) {
          response.json().then(json => {
            resolve({
              json: json,
              status: response.status,
              objectId: response.headers.get("X-Object-Id")
            } as AjaxResult);
          }).catch(err => {
            resolve({
              json: {},
              status: response.status,
              objectId: response.headers.get("X-Object-Id")
            } as AjaxResult);
          });
        } else {
          reject(new Error("Got status code " + response.status));
        }
      }).catch(err => {
        reject(err);
      });
    });
  }

  static async postData(url: string, data?: any): Promise<AjaxResult> {
    return Ajax.query("POST", url, data);
  }

  static async putData(url: string, data?: any): Promise<AjaxResult> {
    return Ajax.query("PUT", url, data);
  }

  static async saveEntity(e: Entity, url: string): Promise<AjaxResult> {
    if (!url.endsWith("/")) {
      url += "/";
    }
    if (e.id) {
      return Ajax.putData(url + e.id, e.serialize());
    } else {
        return Ajax.postData(url, e.serialize()).then(result => {
          e.id = result.objectId;
          return result;
        });
    }
  }

  static async get(url: string): Promise<AjaxResult> {
    return Ajax.query("GET", url);
  }

  static async delete(url: string): Promise<AjaxResult> {
    return Ajax.query("DELETE", url);
  }
}
