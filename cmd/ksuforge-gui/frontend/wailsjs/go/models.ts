export namespace cleanup {
	
	export class Failure {
	    path: string;
	    error: string;
	
	    static createFrom(source: any = {}) {
	        return new Failure(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.path = source["path"];
	        this.error = source["error"];
	    }
	}
	export class File {
	    path: string;
	    bytes: number;
	    logicalBytes: number;
	    modified: string;
	
	    static createFrom(source: any = {}) {
	        return new File(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.path = source["path"];
	        this.bytes = source["bytes"];
	        this.logicalBytes = source["logicalBytes"];
	        this.modified = source["modified"];
	    }
	}
	export class Item {
	    id: string;
	    kind: string;
	    device: string;
	    version: string;
	    name: string;
	    files: File[];
	    bytes: number;
	    recommended: boolean;
	
	    static createFrom(source: any = {}) {
	        return new Item(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.kind = source["kind"];
	        this.device = source["device"];
	        this.version = source["version"];
	        this.name = source["name"];
	        this.files = this.convertValues(source["files"], File);
	        this.bytes = source["bytes"];
	        this.recommended = source["recommended"];
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
	export class Report {
	    token: string;
	    root: string;
	    items: Item[];
	    totalBytes: number;
	
	    static createFrom(source: any = {}) {
	        return new Report(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.token = source["token"];
	        this.root = source["root"];
	        this.items = this.convertValues(source["items"], Item);
	        this.totalBytes = source["totalBytes"];
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
	export class Result {
	    deletedFiles: number;
	    deletedBytes: number;
	    failed: Failure[];
	    cancelled: boolean;
	
	    static createFrom(source: any = {}) {
	        return new Result(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.deletedFiles = source["deletedFiles"];
	        this.deletedBytes = source["deletedBytes"];
	        this.failed = this.convertValues(source["failed"], Failure);
	        this.cancelled = source["cancelled"];
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

}

export namespace device {
	
	export class Device {
	    Serial: string;
	    Mode: string;
	    Manufacturer: string;
	    Brand: string;
	    Model: string;
	    Codename: string;
	    Android: string;
	    Version: string;
	    Slot: string;
	    Kernel: string;
	    KMI: string;
	    HasInitBoot: boolean;
	    SDK: number;
	
	    static createFrom(source: any = {}) {
	        return new Device(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Serial = source["Serial"];
	        this.Mode = source["Mode"];
	        this.Manufacturer = source["Manufacturer"];
	        this.Brand = source["Brand"];
	        this.Model = source["Model"];
	        this.Codename = source["Codename"];
	        this.Android = source["Android"];
	        this.Version = source["Version"];
	        this.Slot = source["Slot"];
	        this.Kernel = source["Kernel"];
	        this.KMI = source["KMI"];
	        this.HasInitBoot = source["HasInitBoot"];
	        this.SDK = source["SDK"];
	    }
	}

}

export namespace main {
	
	export class InstallRequest {
	    Serial: string;
	    Mode: string;
	    Partition: string;
	    OutputDir: string;
	    Catalog: string;
	    KSUd: string;
	    KMI: string;
	    Kernel: string;
	    Mirrors: string;
	    ROMPath: string;
	    Connections: number;
	    Flash: boolean;
	    NoReboot: boolean;
	    VerifyROM: boolean;
	
	    static createFrom(source: any = {}) {
	        return new InstallRequest(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Serial = source["Serial"];
	        this.Mode = source["Mode"];
	        this.Partition = source["Partition"];
	        this.OutputDir = source["OutputDir"];
	        this.Catalog = source["Catalog"];
	        this.KSUd = source["KSUd"];
	        this.KMI = source["KMI"];
	        this.Kernel = source["Kernel"];
	        this.Mirrors = source["Mirrors"];
	        this.ROMPath = source["ROMPath"];
	        this.Connections = source["Connections"];
	        this.Flash = source["Flash"];
	        this.NoReboot = source["NoReboot"];
	        this.VerifyROM = source["VerifyROM"];
	    }
	}
	export class Settings {
	    outputDir: string;
	    language: string;
	    mode: string;
	    partition: string;
	    autoFlash: boolean;
	    noReboot: boolean;
	    catalogUrl: string;
	    mirrors: string;
	    connections: number;
	    ksud: string;
	    logAutoScroll: boolean;
	
	    static createFrom(source: any = {}) {
	        return new Settings(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.outputDir = source["outputDir"];
	        this.language = source["language"];
	        this.mode = source["mode"];
	        this.partition = source["partition"];
	        this.autoFlash = source["autoFlash"];
	        this.noReboot = source["noReboot"];
	        this.catalogUrl = source["catalogUrl"];
	        this.mirrors = source["mirrors"];
	        this.connections = source["connections"];
	        this.ksud = source["ksud"];
	        this.logAutoScroll = source["logAutoScroll"];
	    }
	}
	export class TaskState {
	    status: string;
	    language?: string;
	    log: string;
	    error?: string;
	    result?: any;
	
	    static createFrom(source: any = {}) {
	        return new TaskState(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.status = source["status"];
	        this.language = source["language"];
	        this.log = source["log"];
	        this.error = source["error"];
	        this.result = source["result"];
	    }
	}

}

