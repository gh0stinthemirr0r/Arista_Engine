export namespace core {
	
	export class APIDefinition {
	    id: string;
	    service: string;
	    method: string;
	    path: string;
	    description: string;
	    params: string[];
	    category?: string;
	    tags?: string[];
	
	    static createFrom(source: any = {}) {
	        return new APIDefinition(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.service = source["service"];
	        this.method = source["method"];
	        this.path = source["path"];
	        this.description = source["description"];
	        this.params = source["params"];
	        this.category = source["category"];
	        this.tags = source["tags"];
	    }
	}
	export class APICatalog {
	    eapi: Record<string, APIDefinition>;
	    cloudvision: Record<string, APIDefinition>;
	    eos_rest: Record<string, APIDefinition>;
	    telemetry: Record<string, APIDefinition>;
	    // Go type: time
	    lastUpdated: any;
	
	    static createFrom(source: any = {}) {
	        return new APICatalog(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.eapi = this.convertValues(source["eapi"], APIDefinition, true);
	        this.cloudvision = this.convertValues(source["cloudvision"], APIDefinition, true);
	        this.eos_rest = this.convertValues(source["eos_rest"], APIDefinition, true);
	        this.telemetry = this.convertValues(source["telemetry"], APIDefinition, true);
	        this.lastUpdated = this.convertValues(source["lastUpdated"], null);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	
	export class APIQueryRecord {
	    id: string;
	    endpointId: string;
	    method: string;
	    path: string;
	    body?: Record<string, any>;
	    status: number;
	    response: Record<string, any>;
	    // Go type: time
	    timestamp: any;
	    elapsedMs: number;
	    error?: string;
	
	    static createFrom(source: any = {}) {
	        return new APIQueryRecord(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.endpointId = source["endpointId"];
	        this.method = source["method"];
	        this.path = source["path"];
	        this.body = source["body"];
	        this.status = source["status"];
	        this.response = source["response"];
	        this.timestamp = this.convertValues(source["timestamp"], null);
	        this.elapsedMs = source["elapsedMs"];
	        this.error = source["error"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class ConnectionTestResult {
	    success: boolean;
	    message: string;
	    statusCode?: number;
	    elapsedMs: number;
	    details?: any;
	
	    static createFrom(source: any = {}) {
	        return new ConnectionTestResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.success = source["success"];
	        this.message = source["message"];
	        this.statusCode = source["statusCode"];
	        this.elapsedMs = source["elapsedMs"];
	        this.details = source["details"];
	    }
	}
	export class DeviceInventory {
	    id: string;
	    name: string;
	    deviceType: string;
	    url: string;
	    username: string;
	    password: string;
	    type: string;
	    status: string;
	    // Go type: time
	    addedAt: any;
	    // Go type: time
	    lastTested: any;
	    testCount: number;
	    successCount: number;
	    notes: string;
	
	    static createFrom(source: any = {}) {
	        return new DeviceInventory(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.name = source["name"];
	        this.deviceType = source["deviceType"];
	        this.url = source["url"];
	        this.username = source["username"];
	        this.password = source["password"];
	        this.type = source["type"];
	        this.status = source["status"];
	        this.addedAt = this.convertValues(source["addedAt"], null);
	        this.lastTested = this.convertValues(source["lastTested"], null);
	        this.testCount = source["testCount"];
	        this.successCount = source["successCount"];
	        this.notes = source["notes"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class Endpoint {
	    id: string;
	    name: string;
	    type: string;
	    url: string;
	    username?: string;
	    password?: string;
	    token?: string;
	    // Go type: time
	    created: any;
	    tags: string[];
	    tlsVerify: boolean;
	    status?: string;
	
	    static createFrom(source: any = {}) {
	        return new Endpoint(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.name = source["name"];
	        this.type = source["type"];
	        this.url = source["url"];
	        this.username = source["username"];
	        this.password = source["password"];
	        this.token = source["token"];
	        this.created = this.convertValues(source["created"], null);
	        this.tags = source["tags"];
	        this.tlsVerify = source["tlsVerify"];
	        this.status = source["status"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class ExplorerRequest {
	    endpointId: string;
	    method: string;
	    path: string;
	    body?: Record<string, any>;
	    timeoutMs?: number;
	
	    static createFrom(source: any = {}) {
	        return new ExplorerRequest(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.endpointId = source["endpointId"];
	        this.method = source["method"];
	        this.path = source["path"];
	        this.body = source["body"];
	        this.timeoutMs = source["timeoutMs"];
	    }
	}
	export class ExplorerResponse {
	    status: number;
	    headers: Record<string, Array<string>>;
	    json?: any;
	    text?: string;
	    elapsedMs: number;
	    endpointId: string;
	    logId: string;
	    error?: string;
	
	    static createFrom(source: any = {}) {
	        return new ExplorerResponse(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.status = source["status"];
	        this.headers = source["headers"];
	        this.json = source["json"];
	        this.text = source["text"];
	        this.elapsedMs = source["elapsedMs"];
	        this.endpointId = source["endpointId"];
	        this.logId = source["logId"];
	        this.error = source["error"];
	    }
	}

}

export namespace netvisor {
	
	export class APIDefinition {
	    id: string;
	    service: string;
	    method: string;
	    path: string;
	    description: string;
	    category: string;
	    tags: string;
	    parameters: string;
	    example: string;
	
	    static createFrom(source: any = {}) {
	        return new APIDefinition(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.service = source["service"];
	        this.method = source["method"];
	        this.path = source["path"];
	        this.description = source["description"];
	        this.category = source["category"];
	        this.tags = source["tags"];
	        this.parameters = source["parameters"];
	        this.example = source["example"];
	    }
	}

}

