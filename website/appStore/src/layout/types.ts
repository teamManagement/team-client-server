/**
 * 顶部导航菜单项
 */
export interface HeaderMenuItemInfo {
  /**
   * id
   */
  id: string;

  /**
   * 路由名称
   */
  name: string;

  /**
   * 菜单名称
   */
  title: string;

  /**
   * 图标
   */
  icon: string;

  /**
   * 子菜单
   */
  children: HeaderMenuItemInfo[];
}

export interface MenuItemInfo {
  id: string;
  name: string;
  title: string;
  icon?: any;
  metaData?: string;
  pid?: string;
  type?: number;
  children?: MenuItemInfo[];
  ops?: MenuItemInfo[];
}

