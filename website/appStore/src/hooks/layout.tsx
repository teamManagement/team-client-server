import { useMemo } from "react";
import { useLocation } from "react-router-dom";



/**
 * 使用菜单当前的key
 * @returns [openKeys, selectKeys]
 */
 export function useMenuCurrentKeys(): [string[], string[]] {
  const location = useLocation();
  return useMemo(() => {
    const openKeys: string[] = [];
    const pathName = location.pathname;
    const pathNameSplit = pathName.split("/");
    if (pathNameSplit && pathNameSplit.length > 1) {
      let tmpStr = "";
      for (let str of pathNameSplit) {
        if (!str) {
          continue;
        }
        tmpStr += "/" + str;
        openKeys.push(tmpStr);
      }
    }

    return [openKeys, [pathName]];
  }, [location.pathname]);
}