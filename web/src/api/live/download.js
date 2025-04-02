import request from "@/utils/request";

// 查询角色列表
export function listDownload(query) {
  return request({
    url: "/media/download/list",
    method: "get",
    params: query,
  });
}

export function listDownloadFromCache() {
  return request({
    url: "/media/download/listCache",
    method: "get",
  });
}

export function delRecord(id) {
  return request({
    url: "/media/download/" + id,
    method: "delete",
  });
}
