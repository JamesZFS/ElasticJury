import Axios from "axios";

export const APIS = {
    PING: '/api/ping',
    SEARCH_CASE_ID: '/api/search',
    GET_CASE_INFO: '/api/case',
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
 * @returns {Promise<Object>}
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
 * @returns {Promise<Object>}
 */
export async function getCaseInfo(ids) {
    let res = await Axios.get(APIS.GET_CASE_INFO, {
        params: {id: ids.join(',')}
    });
    return res.data
}
