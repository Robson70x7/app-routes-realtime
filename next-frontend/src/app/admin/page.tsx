"use client";

import { useEffect, useRef } from "react";
import { useMap } from "../../hooks/useMap";
import { socket } from "@/utils/socket-io";

export function AdminPage() {
  const mapContainerRef = useRef<HTMLDivElement>(null);
  const map = useMap(mapContainerRef);
  useEffect(() => {
    if (!map) {
      return;
    }
    // if (socket.disconnected) {
    //   socket.connect();
    // } else {
    //   socket.offAny();
    // }

    //subscribe to on event
    socket.connect();
    console.log("socker disconected? : ", socket.disconnected);
    socket.on("server:new-points:list", async (data) => {
      console.log(data);
      const { route_id, lat, lng } = data;
      const resp = await fetch(`http://localhost:3001/api/routes/${route_id}`);
      const route = await resp.json();
      console.log("route admin", route);
      if (!map.hasRoute(route_id)) {
        map.addRouteWithIcons({
          routeId: route_id,
          startMarkerOptions: {
            position: route.directions.routes[0].legs[0].start_location,
          },
          endMarkerOptions: {
            position: route.directions.routes[0].legs[0].end_location,
          },
          carMarkerOptions: {
            position: route.directions.routes[0].legs[0].start_location,
          },
        });
      }
      map.moveCar(route_id, { lat, lng });
    });

    return () => {
      socket.disconnect();
    };
  }, [map]);

  return (
    <div style={{ width: "100vw", height: "100vh" }} ref={mapContainerRef} />
  );
}

export default AdminPage;
