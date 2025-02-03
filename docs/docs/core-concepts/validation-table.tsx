// // MyTable.tsx
// import React from "react";
// // import styles from "./MyTable.module.css"; // Optional: CSS module for styling

// // Define the type for a single row in the table
// interface TableRow {
//   schemaType: string;
//   data: string;
//   dest: string;
//   requiredError: string;
//   coercionError: string;
// }

// const MyTable: React.FC = () => {
//   // Define the table data array with the specified type
//   const tableData: TableRow[] = [
//     {
//       schemaType: "Bool()",
//       data: "true",
//       dest: "true",
//       requiredError: "no",
//       coercionError: "no",
//     },
//     {
//       schemaType: "Bool()",
//       data: "false",
//       dest: "false",
//       requiredError: "no",
//       coercionError: "no",
//     },
//     {
//       schemaType: "Bool()",
//       data: "nil",
//       dest: "false",
//       requiredError: "yes",
//       coercionError: "yes",
//     },
//     {
//       schemaType: "Bool()",
//       data: '""',
//       dest: "false",
//       requiredError: "yes",
//       coercionError: "yes",
//     },
//     {
//       schemaType: "Bool()",
//       data: '" "',
//       dest: "false",
//       requiredError: "yes",
//       coercionError: "yes",
//     },
//     {
//       schemaType: "Bool()",
//       data: '"on"',
//       dest: "true",
//       requiredError: "no",
//       coercionError: "no",
//     },
//     {
//       schemaType: "Bool()",
//       data: '"off"',
//       dest: "false",
//       requiredError: "no",
//       coercionError: "no",
//     },
//     {
//       schemaType: "Bool()",
//       data: '"true", "t", "T", "True", "TRUE"',
//       dest: "true",
//       requiredError: "no",
//       coercionError: "no",
//     },
//     {
//       schemaType: "Bool()",
//       data: '"false", "f", "F", "FALSE", "False"',
//       dest: "false",
//       requiredError: "no",
//       coercionError: "no",
//     },
//     {
//       schemaType: "Bool()",
//       data: '"test"',
//       dest: "false",
//       requiredError: "no",
//       coercionError: "yes",
//     },
//     {
//       schemaType: "Bool()",
//       data: "1",
//       dest: "true",
//       requiredError: "no",
//       coercionError: "no",
//     },
//     {
//       schemaType: "Bool()",
//       data: "0",
//       dest: "false",
//       requiredError: "no",
//       coercionError: "no",
//     },
//     {
//       schemaType: "Bool()",
//       data: "123",
//       dest: "false",
//       requiredError: "no",
//       coercionError: "yes",
//     },
//     {
//       schemaType: "String()",
//       data: '""',
//       dest: '""',
//       requiredError: "yes",
//       coercionError: "no",
//     },
//     {
//       schemaType: "String()",
//       data: '" "',
//       dest: '""',
//       requiredError: "yes",
//       coercionError: "no",
//     },
//     {
//       schemaType: "String()",
//       data: "nil",
//       dest: '""',
//       requiredError: "yes",
//       coercionError: "yes",
//     },
//     {
//       schemaType: "String()",
//       data: "any value",
//       dest: 'fmt.Sprintf("%v", value)',
//       requiredError: "no",
//       coercionError: "no",
//     },
//     // Add the rest of your data rows here
//   ];

//   return (
//     <table>
//       <thead>
//         <tr>
//           <th>Schema Type</th>
//           <th>Data</th>
//           <th>Dest</th>
//           <th>Required Error (Zero Value)</th>
//           <th>Coercion Error</th>
//         </tr>
//       </thead>
//       <tbody>
//         {tableData.map((row, index) => (
//           <tr key={index}>
//             <td>{row.schemaType}</td>
//             <td>{row.data}</td>
//             <td>{row.dest}</td>
//             <td>{row.requiredError}</td>
//             <td>{row.coercionError}</td>
//           </tr>
//         ))}
//       </tbody>
//     </table>
//   );
// };

// export default MyTable;

const schemaData = [
  {
    schemaType: "Bool()",
    data: "true",
    requiredError: "no",
  },
  {
    schemaType: "Bool()", 
    data: "false",
    requiredError: "yes",
  },
  {
    schemaType: "String()",
    data: '""',
    requiredError: "yes",
  },
  {
    schemaType: "String()",
    data: "any value",
    requiredError: "no",
  },
  {
    schemaType: "Int()",
    data: 0,
    requiredError: "yes",
  },
  {
    schemaType: "Int()",
    data: 10,
    requiredError: "no",
  },
  {
    schemaType: "Float()",
    data: 0.0,
    requiredError: "yes",
  },
  {
    schemaType: "Float()",
    data: 1.21,
    requiredError: "no",
  },
  {
    schemaType: "Time()",
    data: "time.Time{}",
    requiredError: "yes",
  },
  {
    schemaType: "Time()",
    data: "time.Now()",
    requiredError: "no",
  },
  {
    schemaType: "Slice()",
    data: "[]",
    requiredError: "yes",
  },
  {
    schemaType: "Slice()",
    data: "[1,2,3]",
    requiredError: "no",
  },
  {
    schemaType: "Slice()",
    data: "nil",
    requiredError: "yes",
  },
  {
    schemaType: "Ptr(schema)",
    data: "<zero_value_for_schema>",
    requiredError: "no",
  },
  {
    schemaType: "Ptr(schema)",
    data: "nil",
    requiredError: "yes",
  },
];

export default function MyTable() {
  return (
    <div className="">
      <table>
        <thead>
          <tr>
            <th>Schema Type</th>
            <th>Data</th>
            <th>Required Error (Zero Value)</th>
          </tr>
        </thead>
        <tbody>
          {schemaData.map((row, index) => (
            <tr key={index}>
              <td>{row.schemaType}</td>
              <td>{row.data}</td>
              <td>{row.requiredError}</td>
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  );
}
