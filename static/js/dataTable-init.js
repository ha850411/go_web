function initProductTable()
{
    productTable = $('#productTable').on( 'processing.dt', function ( e, settings, processing ) {
        if(processing) {
            $('#loadingTable').css( 'display', 'table');
            $('#productTable_wrapper').css( 'display', 'none');
        } else {
            $('#productTable_wrapper').css( 'display', 'block');
            $('#loadingTable').css( 'display', 'none');
        }
    }).DataTable({
        language: {
            url: "static/vendor/dataTables/languages/ch-tw.json"  
        },
        serverSide:true,
        ajax: {
            url: "/api/products/getList",
            type: "GET",
        },
        columns:[
            {
                data: 'name', 
                title: "產品名稱",
            },
            {
                data: null,
                title: "數量",
                render: function(data, type, row) {
                    output = '<div style="display:flex; align-items:center;">';
                    output += '<input type="number" class="amountInput form-control" style="width:100px;" data-id="'+row.id+'" value="'+row.amount+'" oninput="ListenInput(event, '+row.id+')">';
                    output += '<button type="button" class="btn btn-success btn-sm mdfBtn_'+row.id+'" style="margin-left:5px" onclick="updateAmount('+row.id+')">修改</button>';
                    output += '<i class="fa fa-spinner fa-spin loading_'+row.id+'" style="margin-left:5px;font-size:20px;display:none;"></i>';
                    output += '<i class="fa fa-solid fa-check statusIcon success_'+row.id+'" style="color:#12af55;margin-left:5px;font-size:20px;display:none;"></i>';
                    output += '<i class="fa fa-solid fa-xmark statusIcon failed_'+row.id+'" style="color:#FF0000;margin-left:5px;font-size:20px;display:none;"></i>';
                    output += '</div>';
                    output += '<div class="msgItem msgItem_'+row.id+'" style="display:none;">message: <span class="msg_'+row.id+'" style="color:#FF0000;"></span></div>';
                    return output;
                }
            },
            {
                data: 'amountNotice', 
                title: "存貨警戒線",
                render: function(data, type, row) {
                    if(data>0) {
                        return data;
                    } else {
                        return "無";
                    }
                }
            },
            {
                data: 'formatTime',
                title: "更新時間"
            },
            {
                data: 'id', 
                title: "",
                render: function(data, type, row) {
                    output = '<div style="display:flex; align-items:center;">';
                    output += '<button type="button" class="btn btn-primary btn-sm" style="margin-right:5px">編輯產品</button>';
                    output += '<button type="button" class="btn btn-danger btn-sm">刪除產品</button>';
                    output += '</div>';
                    return output;
                }
            },
        ],
        columnDefs:[
            {
                targets: 4,
                orderable: false,
                responsivePriority: 2,
            },
            {
                targets: 0,
                responsivePriority: 1,
            }
        ],
    });
    return initProductTable;
}