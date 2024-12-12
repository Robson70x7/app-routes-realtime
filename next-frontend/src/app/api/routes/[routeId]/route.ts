import { NextResponse } from "next/server";

export async function GET(
  req: Request,
  { params }: { params: Promise<{ routeId: string }>}
) {
  const { routeId } = await params;
  console.log("next api");
  const res = await fetch(`http://localhost:3000/routes/${routeId}`, {
    cache: "force-cache",
    next: {
      tags: [`routes-${routeId}`, "routes"],
    },
  });
  const data = await res.json();
  return NextResponse.json(data);
}
