import Axios from "axios";

export const APIS = {
    PING: '/api/ping',
    SEARCH_CASE_ID: '/api/search',
    GET_CASE_INFO: '/api/info',
    GET_CASE_DETAIL: '/api/detail/:id',
    GET_ASSOCIATE: '/api/associate/:field/:word'
}

export async function ping() {
    return (await Axios.get(APIS.PING)).data;
}

/**
 * @param misc{string}
 * @param judges{[string]}
 * @param laws{[string]}
 * @param tags{[string]}
 * @return {[int]}
 */
export async function searchCaseId(misc, judges, laws, tags) {
    let res = await Axios.post(APIS.SEARCH_CASE_ID, {
        misc: misc,
        law: laws.join(','),
        tag: tags.join(','),
        judge: judges.join(',')
    }, {responseType: "arraybuffer"})
    let bytes = new Uint8Array(res.data)
    let ids = []
    for (let i = 0; i < bytes.length / 3; ++ i) {
        let id = bytes[i * 3]
        id += bytes[i * 3 + 1] << 8
        id += bytes[i * 3 + 2] << 16
        ids.push(id)
    }
    return ids
}

/**
 * @param ids{[int]}
 * @returns {Promise<[{id: int, judges: [string], laws: [strings], tags: [strings], detail: string}]>}
 */
export async function getCaseInfo(ids) {
    let res = await Axios.post(APIS.GET_CASE_INFO, {
        id: ids.join(',')
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

/**
 *
 * @param field{string}
 * @param word{string}
 * @return {Promise<[string]>}
 */
export async function getAssociate(field, word) {
    let res = await Axios.get(APIS.GET_ASSOCIATE.replace(':field', field).replace(':word', word))
    return res.data.data
}
