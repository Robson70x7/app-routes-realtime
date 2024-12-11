"use client";

import { useRef } from "react";
import { useMap } from "../../hooks/useMap";

export function MapDriver() {
  const mapContainerRef = useRef<HTMLDivElement>(null);
  useMap(mapContainerRef);

  return <div style={{width: "80vw", height: "100vh"}} ref={mapContainerRef} />;
}