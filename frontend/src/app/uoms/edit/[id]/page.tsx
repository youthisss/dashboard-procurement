"use client";

import { SimpleMasterEdit } from "@/components/master-data/simple-master-pages";

export default function UomEdit() {
  return (
    <SimpleMasterEdit
      title="Edit UOM"
      fieldLabel="UOM Name"
      placeholder="e.g. KG"
      requiredMessage="UOM name is required"
    />
  );
}
