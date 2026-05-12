export const onRequest = async (context: any, next: any) => {
  const response = await next();
  
  // Add Link headers for agent discovery (RFC 8288)
  const linkHeaders = [
    '</sitemap.xml>; rel="sitemap"',
    '</about>; rel="about"',
    '</distros>; rel="contents"',
  ];
  
  response.headers.set('Link', linkHeaders.join(', '));
  
  return response;
};
