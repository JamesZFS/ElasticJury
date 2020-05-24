import Axios from "axios";

export const APIS = {
    PING: '/api/ping',
    SEARCH_CASE_ID: '/api/search',
    GET_CASE_INFO: '/api/info',
    GET_CASE_DETAIL: '/api/detail/:id',
}

export async function ping() {
    return (await Axios.get(APIS.PING)).data;
}

/**
 * @param words{[string]}
 * @param judges{[string]}
 * @param laws{[string]}
 * @param tags{[string]}
 * @param misc{string}
 * @returns {Promise<{count: int, result: [int]}>}
 */
export async function searchCaseId(words, judges, laws, tags, misc) {
    let res = await Axios.post(APIS.SEARCH_CASE_ID, {misc}, {
        params: {
            word: words.join(','),
            judge: judges.join(','),
            law: laws.join(','),
            tag: tags.join(','),
        }
    });
    return res.data
}

/**
 * @param ids{[int]}
 * @returns {Promise<[{id: int, judges: [string], laws: [strings], tags: [strings], detail: string}]>}
 */
export async function getCaseInfo(ids) {
    let res = await Axios.get(APIS.GET_CASE_INFO, {
        params: {id: ids.join(',')}
    });
    res.data.forEach(info => {
        for (let field of ['judges', 'laws', 'tags']) {
            info[field] = info[field].split('#').map(s => s.trim()).filter(s => s.length > 0)
        }
    })
    return res.data
}

/**
 * @param id{int}
 * @returns {Promise<{id: int, judges: [string], laws: [strings], tags: [strings], detail: string, tree: string}>}
 */
export async function getCaseDetail(id) {
    let res = await Axios.get(APIS.GET_CASE_DETAIL.replace(':id', String(id)));
    let data = res.data
    for (let field of ['judges', 'laws', 'tags']) {
        data[field] = data[field].split('#').map(s => s.trim()).filter(s => s.length > 0)
    }
    return data
}
