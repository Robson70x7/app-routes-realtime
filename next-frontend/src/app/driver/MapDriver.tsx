"use client";

import { useEffect, useRef } from "react";
import { useMap } from "@/hooks/useMap";
import { socket } from "@/utils/socket-io";

export type MapDriverProps = {
  route_id: string | null;
  start_location: { lat: number; lng: number } | null;
  end_location: { lat: number; lng: number } | null;
};

export function MapDriver(props: MapDriverProps) {
  const { route_id, start_location, end_location } = props;
  const mapContainerRef = useRef<HTMLDivElement>(null);
  const map = useMap(mapContainerRef);

  useEffect(() => {
    if (!map || !route_id || !start_location || !end_location) {
      return;
    }
    if (socket.disconnected) {
      socket.connect();
    } else {
      socket.offAny();
    }

    socket.on("connect", () => {
      console.log("ws connected...");
      socket.emit("client:new-points", { route_id });
    });

    //subscribe to on event
    socket.on(`server:new-points/${route_id}:list`, (data) => {
      console.log(data);
      const { lat, lng } = data;
      if (!map.hasRoute(route_id)) {
        console.log('add route with icons');
        map.addRouteWithIcons({
          routeId: data.route_id,
          startMarkerOptions: {
            position: start_location,
          },
          endMarkerOptions: {
            position: end_location,
          },
          carMarkerOptions: {
            position: start_location,
          },
        });
      }else{
        console.log('move car');
        map.moveCar(route_id, {lat,lng});
      }
    });
    return () => {
      socket.disconnect();
    }
  }, [route_id, start_location, end_location, map]);

  return (
    <div style={{ width: "80vw", height: "100vh" }} ref={mapContainerRef} />
  );
}
