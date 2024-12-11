"use server";
//import { revalidateTag } from "next/cache";

// eslint-disable-next-line @typescript-eslint/no-explicit-any
export async function createRouteAction(state: any, formData: FormData) {
  const { sourceId, destinationId } = Object.fromEntries(formData);

  const directionResponse = await fetch(
    `http://localhost:3000/directions?originId=${sourceId}&destinationId=${destinationId}`,
    {
      cache: "force-cache",
      next: {
        revalidate: 60 * 60, //1H
      },
    }
  );

  if (!directionResponse.ok) {
    return { error: "Falha para carregar direções" };
  }

  const directionData = await directionResponse.json();
  const startAddress = directionData.routes[0].legs[0].start_address;
  const endAddress = directionData.routes[0].legs[0].end_address;

  const response = await fetch(`http://localhost:3000/routes`, {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify({
      name: startAddress + " - " + endAddress,
      source_id: directionData.request.origin.place_id.replace("place_id:", ""),
      destination_id: directionData.request.destination.place_id.replace(
        "place_id:",
        ""
      ),
    }),
  });

  if (!response.ok) {
    return { error: "Falha ao criar rota" };
  }

  return { success: true };
}
