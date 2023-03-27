# Piece Commitment Computation process

(*This is a work in progress*)
Delta is a deal-making service that enables users to make deals with Storage Providers. It is an application that allows users to upload files to the Filecoin network and get them stored by Storage Providers.

A content goes thru different stages in delta. One of the stages is the piece commitment computation process. 

## Content Preparation
- To start, the user will initiate the process by uploading the desired content onto Delta's platform. Once uploaded, Delta will assign the content to a light node and create a corresponding record in its database. To ensure the security and reliability of the content, Delta will assign a miner and wallet to the content record, and create a piece commitment record in its database. This involves computing the piece commitment of the content, which essentially verifies that the content has not been altered in any way.
- Delta will create a deal proposal parameters record in its database, which outlines the terms and conditions of the proposed deal. Based on this record, Delta will then create a deal proposal record, which will serve as a formal agreement between the user and Delta regarding the storage and access of the content on Delta's platform.

## Piece Commitment computation
- The piece commitment dispatched job processes piece commitments for Delta. The purpose of this function is to generate a piece commitment record in Delta's database for the uploaded content, which is later used in creating storage deals with miners.
- It begins by updating the status of the content to "CONTENT_PIECE_COMPUTING" in the database, indicating that the content is being processed. Then, it decodes the CID (content identifier) of the uploaded content and checks for any errors.
- It then prepares to compute the piece commitment of the content by retrieving the content's data using the CID. If Delta's configuration is set to fast mode, the piece commitment is generated using a CommpService. Otherwise, if the connection mode is import, the piece commitment is generated using the filclient library. If neither of these conditions is true, the function generates the piece commitment using the CommpService.
- Once the piece commitment is generated, it is saved to the database as a piece commitment record along with its CID, size, and status. The status of the content in the database is updated to "CONTENT_PIECE_ASSIGNED" to indicate that a piece commitment has been generated, and the ID of the newly created piece commitment record is associated with the content.
- Finally, a new StorageDealMakerProcessor is created with the LightNode, Content, and PieceCommitment record, and it is added to the job queue to [create storage deals](process-flow-storage-deal.md) with miners.