export namespace main {
	
	export class GovInvoiceIndex {
	    Year: string;
	    Month: string;
	    Branch: string;
	    Raw: string;
	
	    static createFrom(source: any = {}) {
	        return new GovInvoiceIndex(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Year = source["Year"];
	        this.Month = source["Month"];
	        this.Branch = source["Branch"];
	        this.Raw = source["Raw"];
	    }
	}
	export class TaskState {
	    NDone1: number;
	    NDone2: number;
	    NDone3: number;
	    Total: number;
	
	    static createFrom(source: any = {}) {
	        return new TaskState(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.NDone1 = source["NDone1"];
	        this.NDone2 = source["NDone2"];
	        this.NDone3 = source["NDone3"];
	        this.Total = source["Total"];
	    }
	}

}

