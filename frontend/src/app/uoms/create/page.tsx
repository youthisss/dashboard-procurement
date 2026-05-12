"use client";

import { SimpleMasterCreate } from "@/components/master-data/simple-master-pages";

export default function UomCreate() {
  return (
    <SimpleMasterCreate
      title="Create UOM"
      fieldLabel="UOM Name"
      placeholder="e.g. KG"
      requiredMessage="UOM name is required"
    />
  );
}
