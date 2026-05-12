"use client";

import { SimpleMasterCreate } from "@/components/master-data/simple-master-pages";

export default function ProductCreate() {
  return (
    <SimpleMasterCreate
      title="Create Product"
      fieldLabel="Product Name"
      placeholder="e.g. CPO"
      requiredMessage="Product name is required"
    />
  );
}
