

export function RFCToFormat(dataStr: any) {
  if (!dataStr) {
    return '-'
  }
  const date = new Date(dataStr).toJSON();
  const newDate = new Date(+new Date(date) + 8 * 3600 * 1000).toISOString().replace(/T/g, ' ').replace(/\.[\d]{3}Z/, '')
  return newDate
}