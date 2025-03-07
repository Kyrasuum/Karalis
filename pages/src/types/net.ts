export const call = (f:any) => typeof f === 'function'? f(): f;

export const parallel = (...fns:any[])  => () => Promise.all(fns.map(call));

export const series = (...fns:any[]) => async () => {
	let res = [];

	for (let f of fns)
		res.push(await call(f));

	return res;
};

export function until(conditionFunction: any) {
	const poll = (resolve: any) => {
		if(conditionFunction()) resolve();
		else setTimeout(_ => poll(resolve), 400);
	}
	return new Promise(poll);
}

export function parseQuery(queryString) {
    let query = {};
    let pairs = (queryString[0] === '?' ? queryString.substr(1) : queryString).split('&');
    for (let pair of pairs) {
        let split_pair = pair.split('=');
        query[decodeURIComponent(split_pair[0])] = decodeURIComponent(split_pair[1] || '');
    }
    return query;
}
