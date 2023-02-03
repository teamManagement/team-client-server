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


/**
 * 用户信息
 */
export interface UserInfo {
  /**
   * 用户ID
   */
  id: string
  /**
   * 用户真实姓名
   */
  name: string
  /**
   * 用户名
   */
  username: string
  /**
   * 身份识别号
   */
  idCode: string
  /**
   * 手机号
   */
  phone: string
  /**
   * 邮箱
   */
  email: string
  /**
   * 头像
   */
  icon: string

  /**
   * 是否为管理员
   */
  isAppStoreManager:boolean
}

