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
 * @returns {Promise<Object>}
 */
export async function searchCaseId(words) {
    let res = await Axios.get(APIS.SEARCH_CASE_ID, {
        params: {word: words.join(',')}
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
