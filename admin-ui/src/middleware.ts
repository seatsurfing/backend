import { NextRequest, NextResponse } from 'next/server';
 
const PUBLIC_FILE = /\.(.*)$/;
 
export async function middleware(req: NextRequest) {
  if (
    req.nextUrl.pathname.startsWith('/_next') ||
    req.nextUrl.pathname.startsWith('/admin/_next') ||
    req.nextUrl.pathname.includes('/api/') ||
    req.nextUrl.pathname.includes('/admin/api/') ||
    PUBLIC_FILE.test(req.nextUrl.pathname)
  ) {
    return;
  }
 
  if (req.nextUrl.locale === 'default') {
    const locale = req.cookies.get('NEXT_LOCALE')?.value || 'en';
     // https://github.com/seatsurfing/backend/issues/166
    const hostname = process.env.FRONTEND_URL ? process.env.FRONTEND_URL : req.url;
 
    return NextResponse.redirect(
      new URL(
        `/admin/${locale}${req.nextUrl.pathname}${req.nextUrl.search}`,
        hostname,
      ),
    );
  }
}
