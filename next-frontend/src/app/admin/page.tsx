"use client";

import { useRef } from "react";
import { useMap } from "../../hooks/useMap";

export function AdminPage() {
  const mapContainerRef = useRef<HTMLDivElement>(null);
  useMap(mapContainerRef);

  return <div style={{width: "100vw", height: "100vh"}} ref={mapContainerRef} />;
}

export default AdminPage;